import java.util.*;

public class Solution {
  public static void main(String[] args) {
    Scanner in = new Scanner(System.in);
    int x = in.nextInt();
    List<Integer> l = new ArrayList<>();
    for (int i = 0; i < x; i++) {
      l.add(in.nextInt());
    }
    for (int i = 0; i < x; i++) {
      System.out.println(l.get(i));
    }
    System.out.println(1);
  }

  public int adder(int x, int y) {
    return x + y;
  }
}
