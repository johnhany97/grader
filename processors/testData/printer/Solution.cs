using System;

namespace Sample
{
  class Test
  {
    public static void Main(string[] args)
    {
      int n = convertFromStdIn();
      for (int i = 0; i < n; i++)
      {
        int x = convertFromStdIn();
        Console.WriteLine(x + 1);
      }
    }
    private static int convertFromStdIn()
    {
      string nString;
      nString = Console.ReadLine();
      return Int32.Parse(nString);
    }
  }
}
