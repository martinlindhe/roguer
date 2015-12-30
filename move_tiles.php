<?php

// one-shot for moving some files

/*
cli question: i have a folder with files named 0.png, 1.png, ...
299.png  and i want to move image 0-79.png to another place. now im
too lazy to write a script, this should be possible with a simple command, right?

i can use grep, but seems number-range matching is not-a-thing :tm: in regexp,
at least according to so: https://stackoverflow.com/q/676467(edited)
*/

$inDir = "resources/assets/tilesets/oddball/ground";

$ifOutDir   = "resources/assets/tilesets/oddball/ground/8x4";
$elseOutDir = "resources/assets/tilesets/oddball/ground/8x12";

for ($i = 0; $i < 500; $i++) {
    $f = $inDir."/".sprintf("%03d", $i).".png";
    if (file_exists($f)) {

        $out = $i <= 79 ? $ifOutDir : $elseOutDir;
        if (!is_dir($out)) {
            mkdir($out);
        }

        $out .= "/".sprintf("%03d", $i).".png";
        rename($f, $out);
        echo "Renamed ".$f." to ".$out."\n";
    }
}
