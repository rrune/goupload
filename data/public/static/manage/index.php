<?php include __DIR__ . '/../helper/checkLogin.php'; ?>
<?php
if ($root != true) {
  header("location:../");
  exit;
}

function getBool($n)
{
  if ($n) {
    return "true";
  } else {
    return "false";
  }
}
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
      <h3>FILES</h3>
      <div class="box-content">
        <form action="filterUser.php" method="post">
          Filter per user: <br>
          <select name="username" required>
            <?php
            foreach ($json as &$user) {
              echo "<option value=" . $user['username'] . ">" . $user['username'] . "</option>";
            }
            ?>
          </select>
          <input type="submit" name="submit" value="Filter">
        </form>
        <?php
        $data = file_get_contents(__DIR__ . '/../../files.json', true);
        $json = json_decode($data, true);
        foreach ($json as &$file) {
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
    <div class="box">
      <h3>LOGINS</h3>
      <div class="box-content">
        <?php
        $data = file_get_contents(__DIR__ . '/../../logins.json', true);
        $json = json_decode($data, true);
        foreach ($json as &$login) {
          echo "<div class='entry'>";
          echo "<b>Username</b>: " . $login["username"] . "<br>";
          echo "<b>Password</b>: " . $login["password"] . "<br>";
          echo "<b>Root</b>: " . getBool($login["root"]) . "<br>";
          echo "<b>Blind</b>: " . getBool($login["blind"]) . "<br>";
          echo "<b>Onetime</b>: " . getBool($login["onetime"]) . "<br>";
          echo "<a href='https://files.qnd.be/manage/removeLogin.php?" . $login["username"] . "'>Remove User</a>";
          echo "</div>";
        }
        ?>
      </div>
    </div>
    <div class="box">
      <h3>CREATE USER</h3>
      <div class="box-content">
        <b>Root</b>: Can see and use this page<br>
        <b>Blind</b>: Can upload files blind, different directory and not in database<br>
        <b>Onetime</b>: User will be removed after uploading a file<br>
        <br>
        <form class="addUser" action="addUser.php" method="post">
          <input type="text" name="username" placeholder="Username" required>
          <input type="text" name="password" placeholder="Password" required>
          <div>
            <input type="checkbox" name="root" value="root" id="root">
            <label for="root">Root</label>
            <input type="checkbox" name="blind" value="blind" id="blind">
            <label for="blind">Blind</label>
          </div>
          <div>
            <input type="checkbox" name="onetime" value="onetime" id="onetime">
            <label for="onetime">Onetime</label>
          </div>
          <input type="submit" name="submit" value="Create">
        </form>
      </div>
    </div>
  </div>

  <div class="sidebar">
    <a href="../logout.php">Logout</a>
    <a href="../">Home</a>
  </div>
</body>

</html>
