# Bake Shop

Bake Shop is an easy way to _dangerously_ and _recklessly_ do subscription payments over Lightning.

All is explained in the [writeup](https://github.com/futurepaul/bakeshop/blob/master/WRITEUP.md)

## Dev setup 

### Frontend

Requires node.js and I'm using [pnpm as well](https://pnpm.io/installation) but should work just fine with npm.

- cd into the `frontend` folder
- `pnpm install`
- `pnpm run dev`
- visit `http://localhost:3000`


### Backend

Replace the macaroon/tls/url with the right values for you

```
cd backend
go get ./...
go run *.go serve --lnd.macaroon /Users/a/.polar/networks/1/volumes/lnd/dave/data/chain/bitcoin/regtest/admin.macaroon --lnd.tls /Users/a/.polar/networks/1/volumes/lnd/dave/tls.cert --lnd.url 127.0.0.1:10004
```


#### REST API's


Initial REST implementation


```
curl --location --request POST 'localhost:8080/bake' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "test",
    "amount": 1000,
    "interval": 100000,
    "times": 10
}'


Example:

{
    "id": "030a102842bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e12047265616412057772697465",   
}
```


```
curl --location --request POST 'localhost:8080/cancel' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "030a102842bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e12047265616412057772697465"
}'
```


```
curl --location --request GET 'localhost:8080/list'


Example:
{
    "items": [
        {
            "id": "030a102942bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e12047265616412057772697465",
            "macaroon": "0201036c6e64022f030a102942bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e120472656164120577726974650002326c6e642d637573746f6d20737562736372696265206d733a313030303020616d6f756e743a313233342074696d65733a3530000006203ca6e51de2c4218c1f3bb115ab94612c8fb0d0927e615eab9deddfd198091eb6",
            "status": "active"
        },
        {
            "id": "030a102a42bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e12047265616412057772697465",
            "macaroon": "0201036c6e64022f030a102a42bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e120472656164120577726974650002326c6e642d637573746f6d20737562736372696265206d733a313030303020616d6f756e743a313233342074696d65733a353000000620fee7c42027d3a66a1446f1d53491f230730f0b23e46d063e078ad5360a29b6ff",
            "status": "active"
        }
    ]
}
```


```
curl --location --request GET 'localhost:8080/details/030a102842bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e12047265616412057772697465'


Example:
{
    "id": "030a102942bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e12047265616412057772697465",
    "macaroon": "0201036c6e64022f030a102942bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e120472656164120577726974650002326c6e642d637573746f6d20737562736372696265206d733a313030303020616d6f756e743a313233342074696d65733a3530000006203ca6e51de2c4218c1f3bb115ab94612c8fb0d0927e615eab9deddfd198091eb6",
    "lndConnect": "lndconnect://test-for-btc-pay.m.staging.voltageapp.io:10009?macaroon=0201036c6e64022f030a102842bc139d31575a547f76c9fec5fba91201301a170a086f6666636861696e120472656164120577726974650002326c6e642d637573746f6d20737562736372696265206d733a313030303020616d6f756e743a313233342074696d65733a353000000620fa17f02688b2bd125275f9dc35415dba475230aac6534def75ce604293a898b0",
    "status": "active"
}
```

### Subscriber

An initial subscriber process has been created to do the pull-based subscription payments once the macaroon has been created by the payer. 

Running: 

Pass it a macaroonHex that has been created with the above bakery API, and the subscriber pubkey where the payment should go. Make sure `accept-keysend` is running on the receiving node.

```
go run *.go subscriber --lnd.tls /Users/a/.polar/networks/1/volumes/lnd/dave/tls.cert --lnd.url 127.0.0.1:10004 --lnd.macaroonHex 0201036c6e64022f030a1003db61c8cdb8787e805ca945cb43fe2b1201301a170a086f6666636861696e120472656164120577726974650002336c6e642d637573746f6d20737562736372696265206d733a31303030303020616d6f756e743a313030302074696d65733a31300000062032d83f97b849dcc98612e0a384e878c7c97d12bb93b583bff96a9195395a00a1 --subscriber.pubkey 02b691fdb9a2526b237bc7ea7205279455b95b87e3b3bb35468d50e6f83a1a1d5e
```
