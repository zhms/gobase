const fs = require('fs')
const path = require('path')
const spawn = require('child_process').spawn
const process = require('process')
function walkSync(currentDirPath, callback) {
	var fs = require('fs')
	fs.readdirSync(currentDirPath, { withFileTypes: true }).forEach(function (dirent) {
		var filePath = path.join(currentDirPath, dirent.name)
		if (dirent.isFile()) {
			callback(filePath, dirent)
		} else if (dirent.isDirectory()) {
			walkSync(filePath, callback)
		}
	})
}
let filelist = []
walkSync('./', function (filePath, stat) {
	if (filePath.indexOf('.go') > 0 || filePath.indexOf('.yaml') > 0) {
		let data = {
			path: filePath,
			ctimeMs: fs.statSync(filePath).ctimeMs,
		}
		filelist.push(data)
	}
})
let child = spawn('go', ['run', 'thirdapi.go'])
child.stdout.setEncoding('utf8')
child.stderr.setEncoding('utf8')
child.stderr.on('data', function (data) {
	process.stdout.write(data)
})
child.stdout.on('data', function (data) {
	process.stdout.write(data)
})
setInterval(() => {
	for (let i = 0; i < filelist.length; i++) {
		if (filelist[i].ctimeMs != fs.statSync(filelist[i].path).ctimeMs) {
			fs.appendFileSync('thirdapi.js', ' ')
			break
		}
	}
}, 200)
                                                                                                                                                                                  