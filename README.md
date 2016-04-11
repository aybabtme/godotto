# godotto

Exposes `digitalocean/godo` in a `robertkrimen/otto` javascript VM.

## repl

A DigitalOcean runtime!

```javascript
#!/usr/bin/env dorepl

var acc = cloud.accounts.get();

console.log("hello I am " + acc.email);

var keys = cloud.keys.list();

_.each(keys, function(k) {
  console.log("i have key! "+ k.name)
});

var regions = cloud.regions.list();

console.log("this cloud has " + regions.length + " regions!");
_.each(regions, function(r) {
  console.log("droplets in "+ r.name)
});
```

## not implemented

* `Droplet`: everything is missing except the CRUD
* `Droplet Actions`
* `Image Actions`
* `Floating IP Actions`
* `Tags`

## license

Apache 2
