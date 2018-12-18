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
const crypto = require('crypto');
const password = 'Password used to generate key';
// Key length is dependent on the algorithm. In this case for aes192, it is
// 24 bytes (192 bits).
// Use async `crypto.scrypt()` instead.
const key = crypto.scryptSync(password, 'salt', 16);
var message = "hello"

var m = cipher(key,message)
console.log(m);

console.log(decipher(key,m));

app.use(express.static('public'));

// http://expressjs.com/en/starter/basic-routing.html
app.get('/one', function(request, response) {
  console.log((request.rawHeaders));
  const readStream = fs.createReadStream('./img.jpg');
	readStream.on('data', (chunk) =>{
	  response.send(cipher(key,chunk.toString('base64')));
	});
});

function cipher(key, message) {
	var iv = crypto.randomBytes(16);
	iv = iv.toString('hex').slice(0, 16);
	var mykey = crypto.createCipheriv('aes-128-cbc', key, iv);
	var mystr = mykey.update(message, 'utf8', 'hex')
	mystr += mykey.final('hex');
	return iv+":"+mystr;
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
// listen for requests :)
const listener = app.listen(3000, function() {
    console.log('Listening to port:  ' + 3000);
});