var toggle = 0
currId = 0
function tableclick(id) {
  $("#transaction-logs").remove()
  if (toggle && currId===id) {
    toggle = 0
    return
  }

  toggle = 1
  currId = id
  console.log(id)
  var html;
  $.ajax({
    url: "/get-logs/" + id,
    context: document.body,
    success: function (done) {
      if (done) {
        html = `
          <table class="table table-condensed" style="font-size:10px">
            <thead>
              <tr>
                <th>Id</th>
                <th>Time</th>
                <th>Category</th>
                <th>Message</th>
                <th>Raw</th>
              </tr> 
            </thead>
            <tbody>`
        for (let i = 0; i < done.length; i++) {
          const el = done[i]
          html += `
                  <tr>
                    <td>${el.Id}</td>
                    <td
                    data-bs-toggle="tooltip"
                    data-bs-placement="top"
                    data-bs-custom-class="custom-tooltip"
                    data-bs-title="${el.CreatedAt}"
                    >${el.Ago}</td>
                    <td>${el.Category}</td>
                    <td>${el.Message}</td>
                    <td>${el.Raw}</td>
                  </tr>`
        }
        html += `
            </tbody>
          </table>`


        $('#transaction-row-' + id).after(`<tr id="transaction-logs"><td colspan="6">${html}</td></tr>`);

        var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
        var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl)
        })

      } else {
        html = "no logs found"
      }
    },
    error: function (err) {
      console.error(err)
    }
  });

}

function getStatusStr() {
  var status = ""
  const f = $('#finished').is(":checked")
  const ru = $('#running').is(":checked")
  const se = $('#syserr').is(":checked")
  const nm = $('#monomey').is(":checked")

  status += (f == true) ? "'FINISH','MFINISH'," : ""
  status += (ru == true) ? "'RUNNING','MANUAL','INC_RUNNING','ERR_FINISH'," : ""
  status += (se == true) ? "'ERROR','SYSTEM_BUSY','GENERAL_ERR'," : ""
  status += (nm == true) ? "'CANNOT_BORROW','NO_BALANCE','GENERAL_ERR','ERR_FINISH'," : ""

  return status.replace(/,*$/, "");
}

function search() {
  search_text = $("#search_box").val()
  const status_text = getStatusStr()
  const newLink = `/all-trades?page=${curr_page}&search=${search_text}&status=${status_text}`
  window.location.replace(newLink);
}

$('#search_box').keypress(function (e) {
  if (e.which == 13)
    search()
});