// server.js
// where your node app starts

// init project
const express = require('express');
const app = express();

app.use(function(req, res, next) {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

const fs = require('fs');
const aes = require('aes-js');
const crypto = require('crypto');
const password = "banana";
const salt = "salt";
console.log(encode(password));
console.log(encode(salt));
var key = crypto.pbkdf2Sync(encode(password), encode(salt), 100, 16, 'sha512');
console.log(key)
app.use(express.static('public'));

// http://expressjs.com/en/starter/basic-routing.html
app.get('/one', function(request, response) {
  //console.log((request.rawHeaders));
  const readStream = fs.createReadStream('./img.jpg');
	readStream.on('data', (chunk) =>{

	  	response.send(cipher(key,padRight(chunk.toString("base64"))));
	});
});
app.get('*', function(req, res){
  res.send('page not found <a href="one">try this</a>', 404);
});
function padRight(message){
	m = message
	for(var i = 0; i < (16-(message.length%16));i++){
		m += " "
	}
	return m
}
function cipher(key, message) {
	var iv = crypto.randomFillSync(new Uint8Array(16)); 
	console.log(iv)
	var aesCbc = new aes.ModeOfOperation.cbc(key, Array.from(iv));
    return aes.utils.hex.fromBytes(iv)+":"+aes.utils.hex.fromBytes(aesCbc.encrypt(aes.utils.utf8.toBytes(message)));
}
function decipher(key, message) {
	var iv = message.split(":")[0];
	console.log(iv.length)
	var message= message.split(":")[1];
	var mykey = crypto.createDecipheriv('aes-128-cbc', key, iv);
	var mystr = mykey.update(message, 'hex', 'utf8')
	mystr += mykey.final('utf8');
	return mystr;
}
function encode(e) {
	var a = [];
	for(var i = 0; i < e.length; i++) a.push(e.charCodeAt(i));
	return new Uint8Array(a);
}
function bytesToHexString(bytes)
{
    if (!bytes)
        return null;

    bytes = new Uint8Array(bytes);
    var hexBytes = [];

    for (var i = 0; i < bytes.length; ++i) {
        var byteString = bytes[i].toString(16);
        if (byteString.length < 2)
            byteString = "0" + byteString;
        hexBytes.push(byteString);
    }

    return hexBytes.join("");
}
// listen for requests :)
const listener = app.listen(3000, function() {
    console.log('Listening to port:  ' + 3000);
});
function stringToUint(string) {

    var charList = string.split(''),
        uintArray = [];
    for (var i = 0; i < charList.length; i++) {
        uintArray.push(charList[i].charCodeAt(0));
    }
    return new Uint8Array(uintArray);
}