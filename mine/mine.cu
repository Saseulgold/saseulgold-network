#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <cuda_runtime.h>

// SHA-256 상수 K 테이블
__device__ __constant__ uint32_t k[64] = {
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5,
    0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
    0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
    0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
    0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc,
    0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7,
    0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
    0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
    0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
    0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3,
    0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5,
    0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
    0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
    0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2
};

// 비트 회전 함수
__device__ inline uint32_t rotr(uint32_t x, uint32_t n) {
    return (x >> n) | (x << (32 - n));
}

// 문자열을 숫자로 변환 (단순히 0-9 숫자만 처리)
__device__ void uint_to_str(uint64_t num, char *str, int max_len) {
    int i = max_len - 1;
    str[i] = '\0';
    i--;
    if (num == 0) {
        str[i] = '0';
        i--;
    }
    while (num > 0 && i >= 0) {
        str[i] = '0' + (num % 10);
        num /= 10;
        i--;
    }
    // Shift the string to the start
    int start = i + 1;
    int j = 0;
    while (str[start] != '\0') {
        str[j++] = str[start++];
    }
    str[j] = '\0';
}

// SHA-256 패딩 처리
__device__ int sha256_pad(const char *input, int input_len, uint8_t *padded, int padded_size) {
    // 현재 구현은 간단히 64바이트 메시지만 처리 (길이가 55바이트 이하일 때)
    if (input_len > 55) return -1;

    // 복사
    for(int i=0; i<input_len; i++) {
        padded[i] = input[i];
    }
    // 패딩 시작
    padded[input_len] = 0x80;
    for(int i=input_len+1; i<56; i++) {
        padded[i] = 0x00;
    }
    // 메시지 길이 (비트 단위)
    uint64_t bit_len = input_len * 8;
    padded[56] = (bit_len >> 56) & 0xFF;
    padded[57] = (bit_len >> 48) & 0xFF;
    padded[58] = (bit_len >> 40) & 0xFF;
    padded[59] = (bit_len >> 32) & 0xFF;
    padded[60] = (bit_len >> 24) & 0xFF;
    padded[61] = (bit_len >> 16) & 0xFF;
    padded[62] = (bit_len >> 8) & 0xFF;
    padded[63] = bit_len & 0xFF;

    return 0;
}

// SHA-256 변환 함수
__device__ void sha256_transform(const uint8_t *data, uint32_t *state) {
    uint32_t w[64];
    // 메시지 파싱 및 스케줄 생성
    for (int i = 0; i < 16; ++i) {
        w[i] = (data[i * 4] << 24) |
               (data[i * 4 + 1] << 16) |
               (data[i * 4 + 2] << 8) |
               (data[i * 4 + 3]);
    }

    for (int i = 16; i < 64; ++i) {
        uint32_t s0 = rotr(w[i - 15], 7) ^ rotr(w[i - 15], 18) ^ (w[i - 15] >> 3);
        uint32_t s1 = rotr(w[i - 2], 17) ^ rotr(w[i - 2], 19) ^ (w[i - 2] >> 10);
        w[i] = w[i - 16] + s0 + w[i - 7] + s1;
    }

    uint32_t a = state[0], b = state[1], c = state[2], d = state[3];
    uint32_t e = state[4], f = state[5], g = state[6], h = state[7];

    for (int i = 0; i < 64; ++i) {
        uint32_t S1 = rotr(e, 6) ^ rotr(e, 11) ^ rotr(e, 25);
        uint32_t ch = (e & f) ^ (~e & g);
        uint32_t temp1 = h + S1 + ch + k[i] + w[i];
        uint32_t S0 = rotr(a, 2) ^ rotr(a, 13) ^ rotr(a, 22);
        uint32_t maj = (a & b) ^ (a & c) ^ (b & c);
        uint32_t temp2 = S0 + maj;
        h = g;
        g = f;
        f = e;
        e = d + temp1;
        d = c;
        c = b;
        b = a;
        a = temp1 + temp2;
    }

    state[0] += a;
    state[1] += b;
    state[2] += c;
    state[3] += d;
    state[4] += e;
    state[5] += f;
    state[6] += g;
    state[7] += h;
}

// 논스 기반 SHA-256 CUDA 커널
__global__ void sha256_nonce_kernel(const char *seed, int seed_len, uint32_t *output, uint64_t num_threads, uint64_t start_nonce, uint64_t *nonce_found) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx >= num_threads) return;

    uint64_t nonce = start_nonce + (uint64_t)idx;

    // 현재 nonce_found가 아직 초기값(0xFFFFFFFFFFFFFFFF)인지 확인
    if (atomicCAS((unsigned long long int*)nonce_found, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF) == 0xFFFFFFFFFFFFFFFF) {
        // 문자열 결합: seed + "," + nonce + ","
        char nonce_str[25];
        uint_to_str(nonce, nonce_str, 25);

        // ',' 삽입
        int comma_pos = seed_len;
        char full_input[64] = {0};
        for(int i=0; i<seed_len; i++) {
            full_input[i] = seed[i];
        }
        full_input[comma_pos++] = ',';
        // nonce_str 뒤에 ',' 추가
        int j = 0;
        while(nonce_str[j] != '\0' && comma_pos < 63) {
            full_input[comma_pos++] = nonce_str[j++];
        }
        full_input[comma_pos++] = ',';
        full_input[comma_pos] = '\0';
        int full_len = comma_pos; // 실제 입력 길이

        // 패딩
        uint8_t padded[64];
        if (sha256_pad(full_input, full_len, padded, 64) != 0) {
            // 패딩 실패 (메시지 길이 초과)
            return;
        }

        // 초기 해시 상태
        uint32_t state[8] = {
            0x6a09e667,
            0xbb67ae85,
            0x3c6ef372,
            0xa54ff53a,
            0x510e527f,
            0x9b05688c,
            0x1f83d9ab,
            0x5be0cd19
        };

        // SHA-256 변환
        sha256_transform(padded, state);
        // 특정 조건 확인 (예: 해시 값 상위 20비트가 0)
        // if (state[0] < 0x00000000 && state[1] < 0x80000000) { // 예시: 상위 20비트가 0이 되도록
        if (state[0] < 0x00000001) { // 예시: 상위 20비트가 0이 되도록
            // 최초로 찾은 nonce만 기록
            if (atomicCAS((unsigned long long int*)nonce_found, 0xFFFFFFFFFFFFFFFF, nonce) == 0xFFFFFFFFFFFFFFFF) {
                // 해시 결과 저장
                for(int i=0; i<8; i++) {
                    output[i] = state[i];
                }
            }
        }
    }
}

int main(int argc, char *argv[]) {
    if (argc != 2) {
        printf("Usage: %s <epoch>\n", argv[0]);
        return 1;
    }

    const char *epoch = argv[1];
    int epoch_len = strlen(epoch);
    const char seed_prefix[] = "blk-";  // 예시로 seed prefix를 "blk-"로 설정
    int seed_prefix_len = strlen(seed_prefix);

    // seed = seed_prefix + epoch
    char seed[64] = {0};
    for(int i=0; i<seed_prefix_len; i++) {
              seed[i] = seed_prefix[i];
    }
    for(int i=0; i<epoch_len; i++) {
        seed[seed_prefix_len + i] = epoch[i];
    }
    int seed_len = seed_prefix_len + epoch_len;

    const uint64_t num_threads = 1024 * 1024;  // 스레드 수 (예: 1,048,576)
    uint32_t *h_output = (uint32_t *)malloc(8 * sizeof(uint32_t));
    uint64_t h_nonce_found = 0xFFFFFFFFFFFFFFFF;

    char *d_seed;
    uint32_t *d_output;
    uint64_t *d_nonce_found;

    cudaMalloc(&d_seed, seed_len * sizeof(char));
    cudaMalloc(&d_output, 8 * sizeof(uint32_t));
    cudaMalloc(&d_nonce_found, sizeof(uint64_t));

    cudaMemcpy(d_seed, seed, seed_len * sizeof(char), cudaMemcpyHostToDevice);
    cudaMemcpy(d_nonce_found, &h_nonce_found, sizeof(uint64_t), cudaMemcpyHostToDevice);

    int threads_per_block = 256;
    int num_blocks_cuda = (num_threads + threads_per_block - 1) / threads_per_block;

    uint64_t start_nonce = 0;
    bool found = false;

    while(!found) {
        // 현재 nonce_found를 초기값으로 재설정
        h_nonce_found = 0xFFFFFFFFFFFFFFFF;
        cudaMemcpy(d_nonce_found, &h_nonce_found, sizeof(uint64_t), cudaMemcpyHostToDevice);

        // 커널 실행
        sha256_nonce_kernel<<<num_blocks_cuda, threads_per_block>>>(d_seed, seed_len, d_output, num_threads, start_nonce, d_nonce_found);
        cudaDeviceSynchronize();

        // nonce_found 복사
        cudaMemcpy(&h_nonce_found, d_nonce_found, sizeof(uint64_t), cudaMemcpyDeviceToHost);

        if (h_nonce_found != 0xFFFFFFFFFFFFFFFF) {
            // nonce가 발견됨
            printf("Final Nonce found: %llu\n", h_nonce_found);
            // 해당 nonce의 해시값을 호스트로 복사하여 출력
            cudaMemcpy(h_output, d_output, 8 * sizeof(uint32_t), cudaMemcpyDeviceToHost);
            printf("Final Hash: %08x %08x %08x %08x %08x %08x %08x %08x\n",
                   h_output[0], h_output[1], h_output[2], h_output[3],
                   h_output[4], h_output[5], h_output[6], h_output[7]);
            found = true;
        } else {
            // nonce가 발견되지 않음, 시작 nonce를 증가시킴
            start_nonce += num_threads;
            printf("Nonce not found in range %llu to %llu. Trying next range...\n", start_nonce, start_nonce + num_threads - 1);
        }
    }

    cudaFree(d_seed);
    cudaFree(d_output);
    cudaFree(d_nonce_found);

    free(h_output);

    return 0;
}

