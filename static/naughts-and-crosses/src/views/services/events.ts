const axios = require('axios').default;
const { CloudEvent, HTTP } = require('cloudevents');

export async function send(type: string) {
    const event = new CloudEvent({
        source: "web-ui",
        type: type,
        datacontenttype: "text/plain",
        dataschema: "https://d.schema.com/my.json",
        subject: "naughts-and-crosses",
    })
    const message = HTTP.binary(event)

    const result = await axios({
        method: "post",
        url: "http://ponglehub.co.uk/events",
        data: message.body,
        headers: message.headers,
        withCredentials: true,
    })

    console.log(`Result: ${result.status}`)
}

export type EventResponse = {
    type: string,
    data: any,
}

export async function receive(): Promise<EventResponse[]> {
    const result = await axios({
        method: "get",
        url: "http://ponglehub.co.uk/events",
        withCredentials: true,
    })

    if (result.status > 399) {
        throw new Error(`Failed to get events: ${result.status}`)
    }

    if (result.status === 204) {
        return []
    }

    let messages: string[] = result.data.messages;

    return messages.map(m => {
        let parsed = JSON.parse(m);
        return {
            type: parsed.type,
            data: parsed.data,
        };
    });
}