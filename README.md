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


## installation

Automated CI status:
- Linux: [![Build Status](https://travis-ci.org/aybabtme/godotto.svg?branch=master)](https://travis-ci.org/aybabtme/godotto)
- Windows: [![Build status](https://ci.appveyor.com/api/projects/status/7v6iq8fs8rdq3ua2/branch/master?svg=true)](https://ci.appveyor.com/project/aybabtme/godotto/branch/master)

### linux

On debian and children:

```bash
wget https://github.com/aybabtme/godotto/releases/download/0.4/dorepl_Linux_x86_64.deb
sudo dpkg -i dorepl_Linux_x86_64.deb
```

Or pull the binary directly:
```bash
wget -qO- https://github.com/aybabtme/godotto/releases/download/0.4/dorepl_Linux_x86_64.tar.gz | tar xvz
```

### darwin

```bash
brew tap aybabtme/homebrew-tap
brew install dorepl
```

### windows (not tested, ymmv)

unzip this file somewhere, i guess!
```bash
https://github.com/aybabtme/godotto/releases/download/0.4/dorepl_Windows_x86_64.zip
```

## not implemented

* `Droplet`: everything is missing except the CRUD
* `Droplet Actions`
* `Image Actions`
* `Floating IP Actions`
* `Tags`

## license

Apache 2
