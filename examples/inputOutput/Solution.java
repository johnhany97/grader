// import java.util.*;

// public class Solution {
//   public static void main(String[] args) {
//     Scanner in = new Scanner(System.in);
//     int x = in.nextInt();
//     List<Integer> l = new ArrayList<>();
//     for (int i = 0; i < x; i++) {
//       l.add(in.nextInt());
//     }
//     for (int i = 0; i < x; i++) {
//       System.out.println(l.get(i));
//     }
//   }

// // docker ps --format='{{.ID}}' | xargs -n 1 docker inspect -f '{{.ID}} {{.State.Running}} {{.State.StartedAt}}' | awk '$2 == "true" && $3 <= "'$(date -v -5M "+%Y-%m-%dT%H:%M:%S%z")'" { print $1 }' | docker kill

//   public int adder(int x, int y) {
//     return x + y;
//   }
// }

// public class Solution {
//   public static void main(String[] args) {
//     System.out.println("hello");
//     while(true) {
      
//     }
//   }
// }

/**
 * 
 * A very fake solution
 * 
 * Build by John Ayad
 * Copyrights reserved 2019 (c)
 * 
 **/
 
import java.util.*;

public class Solution {
  public static void main(String[] args) {
    Scanner scan = new Scanner(System.in);
    int n = scan.nextInt();
    scan.next();
    Map<String, Integer> map = new HashMap<>();
    int maxSoFar = -1;
    for (int i = 0; i < n; i++) {
      String name = scan.next();
      if (!map.hasKey(name)) {
        map.put(name, 0);
      }
      map.put(name, map.get(name) + 1);
      // Update maxSoFar
      maxSoFar = Math.max(maxSoFar, map.get(name));
    }
    // Filter down Map
    Set<String> filtered = new HashSet<>();
    for (Map.Entry<String, Integer> entry : map.entrySet()) {
      if (entry.value() == maxSoFar) {
        filtered.add(entry.key());
      }
    }
    System.out.println(filtered);
    if (filtered.size() == 1) {
      System.out.println(filtered.iterator().next());
      return;
    }
    // We need to sort...
    List<String> list = new ArrayList<String>(filtered);
    Arrays.sort(list, Comparator.reverseOrder());
    System.out.println(list.get(0));
  }
}
