const axios = require('axios')
const _ = require('lodash')
const moment = require('moment')
const fs = require('fs')

require('console-stamp')(console, {
  format: ':date(yyyy/mm/dd HH:MM:ss.l)'
})

const interval = 30000

let main = async () => {
  let mockFile = loadMock()
  if (!mockFile) console.log('No mock file found')
  console.log(` [${mockFile}] starting mock heartbeat...`)
  await sendHeartbeat(mockFile)
}

let loadMock = () => {
  let args = process.argv.slice(2)
  if (_.isEmpty(args)) return

  let mockFile = args[0]
  return mockFile
}

let sleep = () => {
  return new Promise(resolve => setTimeout(resolve, interval))
}

let sendHeartbeat = async (mockFile) => {
  while (true) {
    await sleep()
    let mockData = JSON.parse(fs.readFileSync(`${__dirname}/mock-data/${mockFile}.json`, 'utf8'))
    console.log(` [${mockFile}] sending heartbeat...`)
    let url = `http://localhost:8801/health`
    await axios.post(url, mockData)
  }
}

main().then(() => console.log('done')).catch((err) => console.error(err))