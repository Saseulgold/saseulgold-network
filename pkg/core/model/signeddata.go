package core

type AttributeMap 		= map[string]Ia
type CachedMap 				= map[string]Ia
type SignedDataObject = map[string]Ia

type SignedData struct {
	data				string
	publicKey		string
	signature		string
	hash 				string
	
	cid 				string
	itype 			string
	timestamp 	string
	attributes 	AttributeMap
	cachedUniv	CachedMap
	cachedLocal	CachedMap
}


func (this SignedData) Obj() SignedDataObject {
	return SignedDataObject{
		"data": this.data,
		"publicKey": this.publicKey,
		"signature": this.signature
	}
}

func (this SignedData) Json() string {
	j, _ := json.Marshal(this.Obj())
	return string(j)
}

func (this SignedData) Size() int {
	return len(this.Json())
}