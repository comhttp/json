package jdb

//// // ReadData appends 'data' path prefix for a database read
//func Read(path, col, coin string) (i interface{}) {
//	err := JDB.Read(path+"/"+col, coin, &i)
//	utl.ErrorLog(err)
//	return
//}
//
//// WriteCoin appends 'coins' path prefix for a database write
//func WriteCoin( path, coin string, v interface{}, d interface{}) {
//	JDB.Write(path+"/coins/", coin, v)
//	JDB.Write(path+"/data/"+coin, "info", d)
//}
//
//// WriteCoin appends 'coins' path prefix for a database write
//func WriteCoinImg(path,coin string, i interface{}) {
//	JDB.Write(path+"/data/"+coin, "logo", i)
//}
//
//// WriteCoin appends 'coins' path prefix for a database write
//func WriteCoinData(path,coin, data string, d interface{})  {
//	JDB.Write(path+"/data/"+coin, data, d)
//}
//
//// WriteExchange appends 'exchanges' path prefix for a databaapp/jormse write
//func WriteExchange(path,slug string, v interface{})  {
//	JDB.Write(path+"/exchanges", slug, v)
//}
//
//// ReadCoins reads in all coin data in and converts to bytes for unmarshalling
//func ReadData(path string,v string) [][]byte {
//	data, err := JDB.ReadAll(path)
//	utl.ErrorLog(err)
//	b := make([][]byte, len(data))
//	for i := range data {
//		b[i] = []byte(data[i])
//	}
//	return b
//}
