const inputTask = document.getElementById("input-box");
const listContainer = document.getElementById("task-list");
const popup = document.getElementById("pop-up");

function addTask() {
    if (inputTask.value === "") {
        openPopup(); // Use the popup instead of alert for empty input
    } else {
        createTaskElement(inputTask.value);
    }
    inputTask.value = "";
    saveData();
}

function createTaskElement(task) {
    let li = document.createElement("li");
    li.innerHTML = task;
    listContainer.appendChild(li);
    let span = document.createElement("span");
    span.innerHTML = "\u00d7";
    li.appendChild(span);
}

inputTask.addEventListener("keydown", function (e) {
    if (e.key === "Enter") {
        addTask(); // Trigger addTask on Enter key press
    }
});

listContainer.addEventListener("click", function (e) {
    if (e.target.tagName === "LI") {
        e.target.classList.toggle("checked");
    } else if (e.target.tagName === "SPAN") {
        e.target.parentElement.remove();
        saveData(); // Save data after removing a task
    }
}, false);

function saveData() {
    localStorage.setItem("data", listContainer.innerHTML);
}

function showTask() {
    listContainer.innerHTML = localStorage.getItem("data");
}

showTask();

function openPopup() {
    popup.classList.add("pop-up-active"); // No auto-close timeout
}

function closePopup() {
    popup.classList.remove("pop-up-active");
}
