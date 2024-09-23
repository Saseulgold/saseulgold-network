
class OpMul:
  def __init__(self, *args):
    self.items = args
    self.res = None

    def cloj(this):
      this.res = self.items[0] * self.items[1]
      return this.res

    self.cloj = cloj

  def eval(self):
    res = []
    print('eval: ', self.items)
    for i in self.items:

      if isinstance(i, OpMul) or isinstance(i, OpAdd):
        _i = i.eval()
      else: 
        _i = i

      res.append(_i)

    self.items = res
    return self.cloj(self)

class OpAdd:
  def __init__(self, *args):
    self.items = args
    self.res = None

    def cloj(this):
      this.res = sum(self.items)
      return this.res

    self.cloj = cloj

  def eval(self):
    res = []
    print('eval: ', self.items)
    for i in self.items:

      if isinstance(i, OpMul) or isinstance(i, OpAdd):
        _i = i.eval()
      else: 
        _i = i

      res.append(_i)

    self.items = res
    return self.cloj(self)

if __name__ == "__main__":
  a = OpAdd(1, OpAdd(2,OpAdd(1, 3)))
  b = OpMul(a, OpAdd(1, 2))
  c = OpMul(b, OpAdd(1,1))

  print(c.eval())