
# TODO


forage:
    later: some other "roam" mode, now we finds item of
        type X and moves straight towards it

    

tree:
    humanoid: able to chop down tree, using an axe
    seed: develop into a new tree
    drops: land nearby tree, not in same spot
    plant trees: require "apple seed", etc...


death:
    die from hunger, thirst, coldness
    when npc dies, corpse should stay. need dead state



cooking place:
    should be usable to cook food

    
hut building:
    if basic needs are satisfied (shelter, cooking and fireplace is nearby),
        decide to build a small hut
       
    

farming:
    find wild fruits and berries
    create garden spots around a hut
    plant seeds in garden, to have them nearby housing





generate island:

    - tweak rolling particle length, default to 8
    - generated islands is ugly



# phaser js game engine

need to build custom version of the phaser js lib, see https://github.com/photonstorm/phaser/issues/1937#issuecomment-127626657


    npm install -g grunt-cli

    git clone https://github.com/photonstorm/phaser/
    cd phaser

    npm install

    grunt build --split true

    find in dist: p2.js, phaser.js, pixi.js
----

    npm install -g yo generator-phaser-browserify
---
