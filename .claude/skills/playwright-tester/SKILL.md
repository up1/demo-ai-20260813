---
name: playwright-tester
description: Geenrate test cases with playwright for your web application
---

## Test folder structure with Page object pattern:
```
tests/
  ├── testcases/
  │   ├── test1.spec.js
  │   ├── test2.spec.js
  │   └── ...
  └── pages/
      ├── page1.js
      ├── page2.js
      └── ... 
```

## Best practices for writing test cases with playwright:
* Alwayes use data_testid attributes in your HTML elements to make it easier to select them in your tests. This helps to avoid brittle selectors that can break when the UI changes.
* Always use the Page Object Model (POM) to organize your test code. This helps to keep your tests maintainable and reusable.
* Use descriptive names for your test cases and page objects. This makes it easier to understand what
