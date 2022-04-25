export function timeSince(time: string): string {
    let elapsedMillis = Date.now() - Date.parse(time);
    
    let seconds = elapsedMillis / 1000;
    if (seconds < 60) {
        return humanise(seconds, 'second');
    }

    let minutes = seconds / 60;
    if (minutes < 60) {
        return humanise(minutes, 'minute');
    }
    
    let hours = minutes / 60;
    if (hours < 24) {
        return humanise(hours, 'hour');
    }

    return humanise(hours / 24, 'day');
}

function humanise(period: number, unit: string) {
    let exact = Math.floor(period);

    return `${exact} ${unit}${exact == 1 ? '' : 's'}`;
}