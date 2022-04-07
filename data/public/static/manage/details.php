<?php include __DIR__ . '/../helper/checkLogin.php'; ?>
<?php
if ($root != true) {
    header("location:../");
    exit;
}

$key = $_SERVER['QUERY_STRING'];
$filename = (string) urldecode($key);
$file = __DIR__ . "/../../uploads/" . $filename;
?>

<!DOCTYPE html>
<html lang="en" dir="ltr">

<head>
    <meta charset="utf-8">
    <title>Manage</title>
    <link rel="stylesheet" href="../css/master.css">
    <link rel="stylesheet" href="./style.css">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="shortcut icon" href="../assets/favicon.ico">
</head>

<body>
    <div class="content">
        <div class="box">
            <h3>DETAILS</h3>
            <div class="box-content">
                <?php
                echo "<b>Filename</b>: " . $filename . "<br>";
                echo "<b>Last modified</b>: " . filemtime($file) . "<br>";
                echo "<b>Size</b>: " . (filesize($file) / 1000) . "kB<br>";
                ?>
            </div>
        </div>
    </div>

    <div class="sidebar">
        <a href="../logout.php">Logout</a>
        <a href="./">Back</a>
    </div>
</body>

</html>
