<?php include __DIR__ . '/../helper/checkLogin.php';?>
<?php
if ($root != true) {
  header("location:../");
  exit;
}

$username = $_POST['username'];

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
      <h3>Files</h3>
      <div class="box-content">
        <form action="./filterUser.php" method="post">
          Filter per user: <br>
          <select name="username" required>
            <?php
              foreach ($json as &$user) {
                echo "<option value=". $user['username'] . ">". $user['username'] . "</option>";
              }
            ?>
          </select>
          <input type="submit" name="submit" value="Filter">
        </form>
        <?php
          $fileData = file_get_contents(__DIR__ . '/../../files.json', true);
          $fileJson = json_decode($fileData, true);
          $filtered = array_values(array_filter($fileJson, function ($var) use ($username) {
              return ($var['author'] == $username);
            }));
            foreach ($filtered as &$file) {
              echo "<div class='entry'>";
              echo "<b>File</b>: " . $file["file"] . "<br>";
              echo "<b>Author</b>: " . $file["author"] . "<br>";
              echo "<b>Timestamp</b>: " . $file["timestamp"] . "<br>";
              echo "<b>Link</b>: " . "<a href='https://files.qnd.be/dl?" . $file["short"] . "'>https://files.qnd.be/dl?" . $file["short"] . "</a><br>";
              echo "<a href='https://files.qnd.be/manage/details.php?" . $file["file"] . "'>Details</a> | ";
              echo "<a href='https://files.qnd.be/manage/moveToBlind.php?" . $file["file"] . "'>Move to blind</a> | ";
              echo "<a href='https://files.qnd.be/manage/removeFile.php?" . $file["file"] . "'>Remove File</a>";
              echo "</div>";
            }
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
