# Queryfilter Example

A minimal and practical example on how to use the Queryfilter library for Go.
This is a minimal API implementation with a single endpoint, backed by a SQLite
database.

In this example we're pretending we're a Tshirt shop and we're offering various tshirts.
Through the afore mentioned endpoint we're giving the user of our API the ability
to filter on price, color and size.

## Getting started

**Getting the project**
```console
git clone github.com:Tmw/queryfilter-example.git
cd queryfilter-example
```

**Starting the server**
```console
make server
```

in a separate terminal window, execute:
```console
make call
```

__Note:__ optionally pipe it to JQ if installed on your system to have a more readable output:
```console
```k
make call | jq
```

And observe a single tshirt being returned out of the 12 present in the database.

**Running the tests**
Alternatively we can try running the test suite that contains a few test cases we can look at:

```console
make test
```


## License
[MIT](./LICENSE)
