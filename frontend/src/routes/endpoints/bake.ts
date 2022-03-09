function fakeUuid(): string {
    return 'xxxx-xxxx-xxx-xxxx'.replace(/[x]/g, (c) => {
        const r = Math.floor(Math.random() * 16);
        return r.toString(16);
    });
}
export async function post({ request }) {
    const body = await request.json();
    console.debug(body)



    if (body.name === "badname") {
        return { status: 400, body: "fail" };
    } else {

        const id = fakeUuid();

        const data = { id }

        return { status: 200, body: JSON.stringify(data) };
    }
}
