import org.junit.Test;
import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;

import static org.junit.Assert.*;

public class SolutionTestRunner {
  @Test
  public void adderWorksWithZero() {
    Solution s = new Solution();
    int actual = s.adder(0, 3);
    assertEquals(3, actual);
  }
  public static void main(String[] args) {
      Result result = JUnitCore.runClasses(Solution.class);
      for (Failure failure : result.getFailures()) {
          System.out.println(failure.toString());
      }
      if (result.wasSuccessful()) {
          System.out.println("All tests passed.");
      }
  }
}
