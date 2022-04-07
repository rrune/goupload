<?php include __DIR__ . '/helper/checkLogin.php'; ?>
<?php include __DIR__ . '/helper/randomString.php'; ?>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="../css/master.css">
<link rel="stylesheet" href="../css/textout.css">
<?php
$target_dir = __DIR__ . "/../uploads/";
$blind_dir = __DIR__ . "/../blind/";
if (isset($_POST['blind']) && $_POST['blind'] == "blind") {
  $target_file = $blind_dir . basename($_FILES["fileToUpload"]["name"]);
} else {
  $target_file = $target_dir . basename($_FILES["fileToUpload"]["name"]);
}
$uploadOk = 1;

// Check if file already exists
if (file_exists($target_file)) {
  echo "Sorry, file already exists.";
  $uploadOk = 0;
  exit;
}

// Check file size
/* if ($_FILES["fileToUpload"]["size"] > 500000) {
  echo "Sorry, your file is too large.";
  $uploadOk = 0;
} */

// Check if $uploadOk is set to 0 by an error
if ($uploadOk == 0) {
  echo "Sorry, your file was not uploaded.";
  echo "<br><br><a class='back' href=" . $_SERVER['HTTP_REFERER'] . ">Back</a>";
  // if everything is ok, try to upload file
} else {
  if (move_uploaded_file($_FILES["fileToUpload"]["tmp_name"], $target_file)) {
    if (isset($_POST['blind']) && $_POST['blind'] == "blind") {
      echo "The file " . htmlspecialchars(basename($_FILES["fileToUpload"]["name"])) . " has been uploaded blind.";
    } else {
      $fileData = file_get_contents(__DIR__ . '/../files.json', true);
      $fileJson = json_decode($fileData, true);

      $random = getRandomString(5);
      while (true) {
        $resources = array_values(array_filter($fileJson, function ($var) use ($random) {
          return ($var['short'] == $random);
        }));
        if (isset($resources[0])) {
          $random = getRandomString(5);
        } else {
          break;
        }
      }

      $author = $_SESSION['UserData']['Username'];
      if ($login[0]['onetime'] == true) {
        $author = $_SESSION['UserData']['Username'] . "(Onetime)";
      }

      $new = array(
        "file" => htmlspecialchars(basename($_FILES["fileToUpload"]["name"])),
        "author" => $author,
        "timestamp" => time(),
        "short" => $random
      );
      array_push($fileJson, $new);
      file_put_contents(__DIR__ . '/../files.json', json_encode($fileJson));

      echo "The file " . htmlspecialchars(basename($_FILES["fileToUpload"]["name"])) . " has been uploaded.";
      echo '<a href="https://files.qnd.be/dl?' . $random . '"> Download Link</a>';
      echo "<br><br><a class='back' href=" . $_SERVER['HTTP_REFERER'] . ">Back</a>";

      //delete onetime account
      if ($login[0]['onetime'] == true) {
        $without = array_values(array_filter($json, function ($var) use ($login) {
          return ($var['username'] != $login[0]['username']);
        }));
        $toWrite = json_encode($without, true);
        $file = fopen(__DIR__ . '/../logins.json', 'w+');
        fwrite($file, $toWrite);
        fclose($file);
      }
    }
  } else {
    echo "Sorry, there was an error uploading your file.";
    echo "<br><br><a class='back' href=" . $_SERVER['HTTP_REFERER'] . ">Back</a>";
  }
}
?>