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
$username = (string) urldecode($key);

if ($username == $_SESSION['UserData']['Username']) {
  echo "You can’t remove your own user";
  echo "<br><br><a class='back' href=" . $_SERVER['HTTP_REFERER'] . ">Back</a>";
  exit;
}

if (count($json) > 1) {
  $resources = array_values(array_filter($json, function ($var) use ($username) {
    return ($var['username'] != $username);
  }));
  file_put_contents(__DIR__ . '/../../logins.json', json_encode($resources));
  echo "Removed User " . $username;
  echo "<br><br><a class='back' href=" . $_SERVER['HTTP_REFERER'] . ">Back</a>";
} else {
  echo "Can’t remove only login";
  echo "<br><br><a class='back' href=" . $_SERVER['HTTP_REFERER'] . ">Back</a>";
}
?>
