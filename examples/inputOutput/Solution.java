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
  }

// docker ps --format='{{.ID}}' | xargs -n 1 docker inspect -f '{{.ID}} {{.State.Running}} {{.State.StartedAt}}' | awk '$2 == "true" && $3 <= "'$(date -v -5M "+%Y-%m-%dT%H:%M:%S%z")'" { print $1 }' | docker kill

  public int adder(int x, int y) {
    return x + y;
  }
}

// public class Solution {
//   public static void main(String[] args) {
//     System.out.println("hello");
//     while(true) {
      
//     }
//   }
// }
