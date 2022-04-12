const axios = require('axios').default;
const { CloudEvent, HTTP } = require('cloudevents');

export class Game {
  async getGames() {
    const event = new CloudEvent({
      source: "web-ui",
      type: "naughts-and-crosses.list-games",
      datacontenttype: "text/plain",
      dataschema: "https://d.schema.com/my.json",
      subject: "cha.json",
    })
    const message = HTTP.structured(event)

    const result = await axios({
      method: "post",
      url: "http://ponglehub.co.uk/events",
      data: message.body,
      headers: message.headers,
      withCredentials: true, // Todo: Try this next time
    })

    console.log(`Result: ${result.status}`)
  }
}