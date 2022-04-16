export type EventResponse = {
    type: string,
    data: any,
}

export type SendEvents = {
    send: (type: string, message: string) => void
    stopper: () => void
}

export async function listen(
    receive: (type: string, message: string) => void,
    stopped: () => void
): Promise<SendEvents> { 
    return new Promise((resolve, reject) => {
        let socket = new WebSocket("ws://ponglehub.co.uk/events")

        socket.onopen = (event: Event) => {
            resolve({
                send: function(type: string, message: string): void {
                    console.log(`Sent message -> ${type}: ${message}`);
                },
                stopper: function() {
                    socket.onclose = () => {}
                    socket.onerror = () => {}
                    socket.close();
                }
            });
        };

        socket.onmessage = (event: Event) => {
            console.log(event)
            receive("type", "message")
        }

        socket.onclose = (event: Event) => {
            console.warn("socket connection closed");
            socket.onclose = () => {}
            socket.onerror = () => {}
            stopped();
        };

        socket.onerror = (error: Event) => {
            console.warn(`Socket error: ${error}`)
            socket.onclose = () => {}
            socket.onerror = () => {}
            reject();
        }
    });   
}