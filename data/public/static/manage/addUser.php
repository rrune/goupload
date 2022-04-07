<?php include __DIR__ . '/../helper/checkLogin.php';?>
<?php include __DIR__ . '/../helper/randomString.php';?>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="../css/master.css">
<link rel="stylesheet" href="../css/textout.css">
<?php
if ($root != true) {
  header("location:../");
  exit;
}

$username = $_POST['username'];
$username = strtolower($username);

$doesAlreadyExist = array_values(array_filter($json, function ($var) use ($username) {
  return ($var['username'] == $username);
}));
if (count($doesAlreadyExist) > 0) {
  echo $_POST['username'] . "already exists";
  echo "<br><br><a class='back' href=". $_SERVER['HTTP_REFERER'] . ">Back</a>";
  exit;
}

$root = false;
$blind = false;
$onetime = false;

if (isset($_POST['root']) && $_POST['root'] == "root") {
    $root = true;
}
if (isset($_POST['blind']) && $_POST['blind'] == "blind") {
    $blind = true;
}
if (isset($_POST['onetime']) && $_POST['onetime'] == "onetime") {
    $onetime = true;
}

$user = array (
    "username" => $username,
    "password" => $_POST['password'],
    "key" => getRandomString(8),
    "root" => $root,
    "blind" => $blind,
    "onetime" => $onetime
);

array_push($json, $user);
file_put_contents(__DIR__ . '/../../logins.json', json_encode($json));

echo "Created User " . $_POST['username'];
echo "<br><br><a class='back' href=". $_SERVER['HTTP_REFERER'] . ">Back</a>";
?>
