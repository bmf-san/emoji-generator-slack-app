<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head> 
<body onLoad="draw()">
<canvas id="target" height="128" width="128"></canvas>
<script>
function draw() {
    var element = document.getElementById("target");
    var context = element.getContext("2d");
    var maxWidth = element.width;

    context.clearRect(0, 0, element.width, element.height);
    context.textAlign = "center";
    context.font = "bold 64px Arial";
    // TODO: 文字列は入力を受け取る。入力は1行もしくは2行とする。1行当たりの文字数に制限はない。
    context.fillText("Hello", element.width*0.5, 56, maxWidth);
    context.fillText("World", element.width*0.5, 120, maxWidth);
}
</script>
</body>
</html>