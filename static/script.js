
  
  function loadLogs(id){
    $.ajax({ url: "/get-logs/" + id,
    context: document.body,
    success: function(done){
      console.log("load available done", done)
      if (done){

      } else {

      }
      $('#spinner').hide()
    }
  });
  }
  function tableclick(id) {
    console.log(id)
    $("#transaction-logs").remove()
    $('#transaction-row-' + id).after(`<tr id="transaction-logs"><td colspan="6">

      <table class="table table-condensed">
        <thead>
          <tr>
            <th>Firstname</th>
            <th>Lastname</th>
            <th>Email</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>John</td>
            <td>Doe</td>
            <td>john@example.com</td>
          </tr>
          <tr>
            <td>Mary</td>
            <td>Moe</td>
            <td>mary@example.com</td>
          </tr>
          <tr>
            <td>July</td>
            <td>Dooley</td>
            <td>july@example.com</td>
          </tr>
        </tbody>
      </table>

    </td></tr>`);
  }

  function populateLinks(page = null) {
    console.log('populating links', page)
    if (!page) {
      page = 21;
      console.log("current page", page)
      if (page == 1) {
        $('#linkprev').addClass('disabled');
      } else {
        $('#linkprev').removeClass('disabled');
        const np = `/?page=${page-1}`;
        console.log(np)
        $('#hlinkprev').attr("href", np)
      }

      for (i=1; i<=3; i++) {
        const np = `/?page=${page+i}`;
        $(`#hlink${i}`).attr("href", np)
        $(`#hlink${i}`).text(page+i)
      }
      $('#hlinkfwd').attr("href", `/?page=${page+4}`)
    }
  }

  $( document ).ready(populateLinks());
