// https://stackoverflow.com/questions/6312993/javascript-seconds-to-time-string-with-format-hhmmss
export default function prettyPrintInterval(secs: number) {
    // const secs = Math.floor(interval / 1000);

    if (typeof secs != 'number') {
        return '0 seconds';
    }

    const days = Math.floor(secs / (60 * 60 * 24));

    const divisor_for_hours = secs % (60 * 60 * 24);
    const hours = Math.floor(divisor_for_hours / (60 * 60));

    const divisor_for_minutes = secs % (60 * 60);
    const minutes = Math.floor(divisor_for_minutes / 60);

    const divisor_for_seconds = divisor_for_minutes % 60;
    const seconds = Math.ceil(divisor_for_seconds);

    if (days > 0) {
        return `${days} days, ${hours} hours, ${minutes} minutes, and ${seconds} seconds`;
    } else if (hours > 0) {
        return `${hours} hours, ${minutes} minutes, and ${seconds} seconds`;
    } else if (minutes > 0) {
        return `${minutes} minutes and ${seconds} seconds`;
    } else {
        return `${seconds} seconds`;
    }
}