import java.util.Scanner;

public class Solution {
  public static void main(String[] args) {
    Scanner scan = new Scanner(System.in);
    String s = scan.nextLine();
    scan.close();
    int count = 0;
    for (Character c : s.toCharArray()) {
      if (c == 'K') {
        count++;
      }
    }
    System.out.println(count);
  }
}
