/** @type {import('@sveltejs/kit').RequestHandler} */
export async function get({ params }) {
    const id = params.uuid
    if (id) {
        const lndconnect = "lndconnect://test-for-btc-pay.m.staging.voltageapp.io:10009?macaroon=AgEDbG5kAvgBAwoQz7yEtz58fUUyb3q_CEkHIxIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaIQoIbWFjYXJvb24SCGdlbmVyYXRlEgRyZWFkEgV3cml0ZRoWCgdtZXNzYWdlEgRyZWFkEgV3cml0ZRoXCghvZmZjaGFpbhIEcmVhZBIFd3JpdGUaFgoHb25jaGFpbhIEcmVhZBIFd3JpdGUaFAoFcGVlcnMSBHJlYWQSBXdyaXRlGhgKBnNpZ25lchIIZ2VuZXJhdGUSBHJlYWQAAAYgBsiUkn1EAT6H1yQ-jVSz0WcRYC5F9lcrNrDP9uQP6Ow"
        return {
            body:
                { id, interval: 3600, name: 'this-is-the-test', times: 21, amount: 21000, lndconnect }
        };
    }

    return {
        status: 404
    };
}