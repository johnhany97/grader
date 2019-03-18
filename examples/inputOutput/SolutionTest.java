import org.junit.Test;
import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;

import static org.junit.Assert.*;

public class SolutionTest {
  @Test
  public void adderWorksWithZero() {
    Solution s = new Solution();
    int actual = s.adder(0, 3);
    assertEquals(3, actual);
  }
}
