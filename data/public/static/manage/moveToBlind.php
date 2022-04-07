<?php include __DIR__ . '/../helper/checkLogin.php'; ?>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="../css/master.css">
<link rel="stylesheet" href="../css/textout.css">
<?php
if ($root != true) {
    header("location:../");
    exit;
}

$key = $_SERVER['QUERY_STRING'];
$filename = (string) urldecode($key);

$fileData = file_get_contents(__DIR__ . '/../../files.json', true);
$fileJson = json_decode($fileData, true);

$resources = array_values(array_filter($fileJson, function ($var) use ($filename) {
    return ($var['file'] != $filename);
}));
file_put_contents(__DIR__ . '/../../files.json', json_encode($resources));

if (!rename(__DIR__ . '/../../uploads/' . $filename, __DIR__ . '/../../blind/' . $filename)) {
    echo ("$filename cannot be moved due to an error");
} else {
    echo ("$filename has been moved");
}
echo "<br><br><a class='back' href=" . $_SERVER['HTTP_REFERER'] . ">Back</a>";
?>
