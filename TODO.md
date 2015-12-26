
# TODO


- freezing Level
    to decide to light a fire
    requires fireplace (if not next to it, move to fireplace)
    put firewood in fireplace
    light fire
    each fire tick, consume energy until fire is burned out
    generate heat. warm up nearby npc (npc can "crouch next to the fire" to get more heat)


- sleeping
    if shelter is nearby, travel to it and then sleep there to get bonus
    
- forage
    move around, dont revisit same place

- trees
    spawn trees over island
    tick tree: chance of dropping tree branches for firewood, chance of dropping seed
    humanoid: able to chop down tree, using an axe
    seed: develop into a new tree





    generated islands is ugly



generate island ui

    - tweak rolling particle length, default to 8




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
