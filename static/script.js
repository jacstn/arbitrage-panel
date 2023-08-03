var toggle = 0;
currId = 0;

const xhr = new XMLHttpRequest();

function addSuccessAlert(msg, id) {
  $("#transaction-row-" + id).after(
    `<tr id="transaction-logs"><td colspan="9">
          <div class="alert alert-success" role="alert">
${msg}</div>
          </td></tr>`
  );
}

function addErrorAlert(msg, id) {
  console.log("adding danger alert");
  $("#transaction-row-" + id).after(
    `<tr id="transaction-logs"><td colspan="9">
          <div class="alert alert-danger" role="alert">
${msg}</div>
          </td></tr>`
  );
}

function onDelayRequestEnd(e) {
  const res = JSON.parse(e.target.responseText);
  console.log(res);
  if (res.status === "err") {
    console.log("error while delaying", res);
  } else if (res.status == "ok") {
    addSuccessAlert("Trade Delayed ", res.data);
  }
}

let selectedToClose = undefined;

function closeTrade(id) {
  const newEl = document.getElementById("close-link-" + id);

  if (selectedToClose !== id) {
    const oldEl = document.getElementById("close-link-" + selectedToClose);
    if (oldEl) {
      oldEl.className = "";
    }
    selectedToClose = id;

    newEl.className = "btn btn-danger btn-sm";
    return false;
  }
  newEl.className = "";
  console.log(id);
  const request = new XMLHttpRequest();
  request.open("GET", "/close-trade/" + id);
  request.send();
  request.addEventListener("loadend", (data) => {
    const { target } = data;
    const { responseText } = target;
    console.log(target, JSON.parse(responseText).msg);
    if (target.status === 502) {
      addErrorAlert(JSON.parse(responseText).msg, id);
      return false;
    }
  });
  return true;
}

function tradeLogs(id) {
  xhr.open("GET", "/get-logs/" + id);
  xhr.send();
  xhr.on;
}

xhr.addEventListener("load", () => {
  if (xhr.status === 200) {
    const resp = JSON.parse(xhr.response);
    if (resp.msg === "trade succesfully closed") {
      console.log(resp);
      const element = document.getElementById("close-link-" + resp.data);
      element.parentNode.removeChild(element);
    }
  } else {
    console.log(xhr.status);
  }
});

function delay(id) {
  console.log("delay reqeust");
  xhr.open("GET", "/delay-trade/" + id);
  xhr.send();
  xhr.addEventListener("loadend", onDelayRequestEnd);
  return false;
}

function tableclick(id) {
  console.log("asdfasdfsad");

  if (currId === id) {
    return;
  }
  $("#transaction-logs").remove();

  currId = id;
  //console.log(id);
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
            <tbody>`;
        for (let i = 0; i < done.length; i++) {
          const el = done[i];
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
                  </tr>`;
        }
        html += `
            </tbody>
          </table>`;

        $("#transaction-row-" + id).after(
          `<tr id="transaction-logs"><td colspan="9">${html}</td></tr>`
        );

        var tooltipTriggerList = [].slice.call(
          document.querySelectorAll('[data-bs-toggle="tooltip"]')
        );
        var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
          return new bootstrap.Tooltip(tooltipTriggerEl);
        });
      } else {
        html = "no logs found";
      }
    },
    error: function (err) {
      console.error(err);
    },
  });
}

function getStatusStr() {
  var status = "";
  const f = $("#finished").is(":checked");
  const ru = $("#running").is(":checked");
  const se = $("#syserr").is(":checked");
  const nm = $("#monomey").is(":checked");

  status += f == true ? "'FINISH','MFINISH'," : "";
  status += ru == true ? "'RUNNING','MANUAL','INC_RUNNING','ERR_FINISH'," : "";
  status += se == true ? "'ERROR','SYSTEM_BUSY','GENERAL_ERR'," : "";
  status +=
    nm == true
      ? "'CANNOT_BORROW','NO_BALANCE','GENERAL_ERR','ERR_FINISH',"
      : "";

  return status.replace(/,*$/, "");
}

function search() {
  search_text = $("#search_box").val();
  const status_text = getStatusStr();
  const newLink = `/all-trades?page=${curr_page}&search=${search_text}&status=${status_text}`;
  window.location.replace(newLink);
}

$("#search_box").keypress(function (e) {
  if (e.which == 13) search();
});
