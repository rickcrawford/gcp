package handlers

import (
	"fmt"
	"net/http"
)

const indexContent = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
    <title>Trie Search</title>

    <!-- Latest compiled and minified CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

    <!-- Optional theme -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

    <!-- Latest compiled and minified JavaScript -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
  </head>
  <body>

      <div class="container">

      <div>
        <h1>Trie Search</h1>

        <form onsubmit="return false;">
            <input id="hero-demo" autofocus type="text" name="q" class="form-control" style="width:500px">
        </form>
      </div>

    </div><!-- /.container -->


   
    <script async src="/script.js?TA_CLIENT_ID=694ce192fcf54e258bbd821e6afdb61e&BATCH_ID=%s"></script>
    <script type="text/javascript">
        window._ta = window._ta || [];
        function typeahead(){_ta.push(arguments)};
        typeahead('TA_CLIENT_ID', '694ce192fcf54e258bbd821e6afdb61e', 'input[name="q"]');
        typeahead('debug', true);
        typeahead('minChars', 1);
        typeahead('onSelect', function(e, term, item){
           console.log(e, term, item);
        });
    </script>


  </body>
</html>

`

func indexHandler(rw http.ResponseWriter, req *http.Request) {

	batch := req.FormValue("BATCH_ID")
	if batch == "" {
		batch = "1"
	}
	fmt.Fprintf(rw, indexContent, batch)
}
