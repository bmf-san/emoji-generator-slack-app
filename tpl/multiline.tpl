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

    element.style.background = {{.BgColor}};
    context.clearRect(0, 0, element.width, element.height);
    context.textAlign = "center";
    context.font = "bold 64px Arial";
    context.fillStyle = {{.Color}};
    // FIXME: Optimal size varies depending on the size, so calculation is required
    context.fillText({{.Line1}}, element.width*0.5, 56, maxWidth);
    context.fillText({{.Line2}}, element.width*0.5, 115, maxWidth);
}
</script>
</body>
</html>