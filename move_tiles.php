<?php

// one-shot for moving some files

/*
cli question: i have a folder with files named 0.png, 1.png, ...
299.png  and i want to move image 0-79.png to another place. now im
too lazy to write a script, this should be possible with a simple command, right?

i can use grep, but seems number-range matching is not-a-thing :tm: in regexp,
at least according to so: https://stackoverflow.com/q/676467(edited)
*/

$inDir = "resources/assets/tilesets/oddball/tiles";

$ifOutDir = "resources/assets/tilesets/oddball/tiles/4x12";
$elseOutDir = "resources/assets/tilesets/oddball/tiles/8x12";

for ($i = 0; $i < 400; $i++) {
    $f = $inDir."/".$i.".png";
    if (file_exists($f)) {

        $out = $i <= 79 ? $ifOutDir : $elseOutDir;
        $out .= "/".$i.".png";
        rename($f, $out);
        echo "Renamed ".$f." to ".$out."\n";
    }
}
