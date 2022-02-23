const axios = require('axios');
const { CloudEvent, HTTP } = require('cloudevents');

export class Game {
  private async getGames() {
    const event = new CloudEvent({
      source: "web-ui",
      type: "type",
      datacontenttype: "text/plain",
      dataschema: "https://d.schema.com/my.json",
      subject: "cha.json",
      data: "my-data",
      extension1: "some extension data",
    })

  }
}