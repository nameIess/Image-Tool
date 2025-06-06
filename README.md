# 📝 To-Do List Application

A simple, responsive, and user-friendly To-Do List web application built with **HTML**, **CSS**, and **JavaScript**. Users can add, complete, and delete tasks — with all tasks saved in the browser via **localStorage** for persistence across sessions.

---

## ✨ Features

- **Add Tasks:** Quickly add tasks via the input field using the **+Task** button or **Enter** key.
- **Complete Tasks:** Click a task to toggle it as completed with a checkmark and strikethrough.
- **Delete Tasks:** Remove tasks with a single click on the `×` icon.
- **Persistent Storage:** Automatically saves tasks to `localStorage`.
- **Input Validation:** Prevents adding empty tasks using a stylish pop-up.
- **Responsive Design:** Optimized layout for both mobile and desktop (max width: 540px).
- **Interactive UI:** Smooth animations, hover effects, and modern gradients.
- **Keyboard Support:** Press **Enter** to add tasks and navigate the interface with ease.

---

## 📁 File Structure

```
/project-root
├── index.html           # Main HTML file with UI and pop-up layout
├── readme.md            # This README file
└── /resource
    ├── script.js        # JavaScript logic for task handling and storage
    └── style.css        # CSS for layout, styling, and animations
```

---

## 🚀 Getting Started

### ✅ Prerequisites

- A modern web browser (e.g., Chrome, Firefox, Edge).
- Optional: Internet connection (for external icon URLs, or host them locally).

### 🔧 Installation

1. **Clone or Download:**
   - Clone the repo or download and extract the ZIP file.
   - Ensure the file structure matches the one shown above.

2. **Run Locally:**
   - Open `index.html` directly in your browser.
   - (Optional) Use a local development server like `live-server` for a better dev experience.

3. **Optional Configuration:**
   - **Fix Broken Image:** Replace or remove the `<img>` in `index.html` (line 10).
   - **Offline Icons:** Host checkbox icons locally and update their URLs in `style.css`.

---

## 🛠 Usage

### ➕ Add Tasks
- Type your task in the "Add your Task" field.
- Click the **+Task** button or press **Enter**.

### ✅ Complete / ❌ Delete Tasks
- **Toggle Complete:** Click on a task to mark/unmark it as completed.
- **Delete:** Click the `×` icon to remove a task.

### ⚠ Input Validation
- If the input is empty, a pop-up appears:
  > _"Task can't be empty. Please enter task before Adding Task."_

- Click **Okay** to close the pop-up and return to the input.

### 💾 Data Persistence
- All tasks are stored using `localStorage`.
- Your list is preserved even after reloading or closing the browser.

---

## 🔍 Development Notes

### 🧱 HTML (`index.html`)
- Contains main UI elements, task list, and validation pop-up.

### 🎨 CSS (`style.css`)
- Uses a gradient background: `linear-gradient(135deg, #153677, #d74072)`.
- Responsive design with animations and hover effects.
- External icons for task state (can be replaced for offline use).

### ⚙ JavaScript (`script.js`)
- Manages tasks: add, complete, delete.
- Stores task list HTML in `localStorage`.
- Handles pop-up display and keyboard events.

---

## 📈 Potential Improvements

### ♿ Accessibility
- Add `aria-label`s and `role="alert"` for better screen reader support.
- Improve keyboard navigation (e.g., `Esc` to close pop-up, `Tab` to access buttons).

### ✏️ UX Features
- Edit tasks by clicking the text.
- Add categories, priority levels, or color tags.
- Include a "Clear All" button.

### 🧹 Pop-Up Enhancements
- Auto-close the pop-up after a delay (`setTimeout`).
- Allow closing pop-up with `Esc`.

### 🛡 Error Handling
- Sanitize localStorage data to avoid rendering malformed HTML.
- Handle long or special character inputs.

### 📱 Mobile Optimization
- Adjust padding, font sizes, and button sizes for touch-friendly interaction.

---

## 🐞 Known Issues

- **Broken Image:** Missing `src` in `<img>` (line 10, `index.html`). Replace or remove.
- **External Icons:** Icons in `style.css` rely on CDN; may not load offline.
- **Manual Pop-Up Closure:** Currently requires clicking **Okay**, which may interrupt user flow.

---

## 🤝 Contributing

Feel free to fork the repo, make changes, and submit pull requests!  
Bug reports and feature suggestions are also welcome.

---

## 📄 License

This project is **unlicensed** and free to use, modify, or distribute for any purpose.

---

## 🙏 Acknowledgments

- Built using vanilla **HTML**, **CSS**, and **JavaScript** for simplicity and educational value.
- Gradient UI inspired by modern minimalistic design trends.