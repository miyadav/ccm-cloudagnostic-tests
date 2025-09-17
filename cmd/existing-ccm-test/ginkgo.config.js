module.exports = {
  // Ginkgo configuration for CCM tests
  // This file configures Ginkgo test runner behavior
  
  // Test timeout
  timeout: "5m",
  
  // Parallel execution
  parallel: false, // Set to true if tests can run in parallel
  
  // Random seed for test execution order
  randomizeAllSpecs: true,
  
  // Fail fast on first failure
  failFast: true,
  
  // Show progress during test execution
  showProgress: true,
  
  // JUnit reporter configuration
  junitReport: {
    enabled: true,
    outputDir: "./test-results",
    filename: "junit.xml"
  }
};
