# godotto

Exposes `digitalocean/godo` in a `robertkrimen/otto` javascript VM.

## repl

A DigitalOcean repl!

```javascript
> cloud.droplets.list();
[]
> cloud.droplets.create({"name":"lol","image":{"slug":"debian-8-x64"},"region":"nyc3","size":"1gb"});
{
  "created_at": "2016-04-11T05:39:19Z",
  "disk": 30,
  "id": 13190234,
  "name": "lol",
  ...
}
> var d = cloud.droplets.get(droplets[0].id);
> d.status;
"active"
> cloud.droplets.delete(d);
> cloud.droplets.list();
[]
```

Or use the REPL as a runtime:

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

## API docs

Insert here.

## not implemented

* `Droplet`: everything is missing except the CRUD
* `Tags`

## license

Apache 2
