
# TODO


forage:
  - later: some other "roam" mode, now we finds item of type X and moves straight towards it



tree:
  - humanoid: able to chop down tree, using an axe
  - seed: develop into a new tree
  - drops: land nearby tree, not in same spot
  - plant trees: require "apple seed", etc...


death:
  - die from hunger, thirst, coldness
  - when npc dies, corpse should stay. need dead state



cooking place:
  - should be usable to cook food


hut:
  - ...


farming:
  - find wild fruits and berries


garden:
  - when we have a hut:
  - create garden spots around a hut
  - plant seeds in garden, to have them nearby housing





generate island:
  - tweak rolling particle length, default to 8
  - generated islands is ugly




# phaser js game engine




  - tool: read *.yml from oddball tileset, use it to pick the proper tiles, and generate new tile maps,
    one for each input folder (so we'll have multiple tile maps still, but mapped to the items we use)




----

NOW:
    - spawning sometimes outside of map. perhaps because spawn data is done before island map has been processed?

    - able to use characters.yml and items.yml on map

    - move ground/8x12 to a sprite layer, rather than ground-tile layer
