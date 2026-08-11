// document_scanner.cs — C# версия

using System;
using System.Collections.Generic;
using System.Drawing;
using System.IO;
using OpenCvSharp;

class Program
{
    static void Main(string[] args)
    {
        if (args.Length == 0)
        {
            Console.WriteLine("Usage: dotnet run <image>");
            return;
        }
        string input = args[0];
        Console.WriteLine("📄 Document Scanner (C#)");
        Console.WriteLine($"📂 Загружено: {input}");
        var start = DateTime.Now;

        using var src = Cv2.ImRead(input);
        if (src.Empty())
        {
            Console.WriteLine("❌ Ошибка загрузки");
            return;
        }
        using var gray = new Mat();
        Cv2.CvtColor(src, gray, ColorConversionCodes.BGR2GRAY);
        using var blurred = new Mat();
        Cv2.GaussianBlur(gray, blurred, new Size(5,5), 0);
        using var edged = new Mat();
        Cv2.Canny(blurred, edged, 75, 200);

        var contours = new Point[][] { };
        var hierarchy = new Mat();
        Cv2.FindContours(edged, out contours, out hierarchy, RetrievalModes.List, ContourApproximationModes.ApproxSimple);

        Point[] screenCnt = null;
        foreach (var cnt in contours)
        {
            var peri = Cv2.ArcLength(cnt, true);
            var approx = Cv2.ApproxPolyDP(cnt, 0.02 * peri, true);
            if (approx.Length == 4)
            {
                screenCnt = approx;
                break;
            }
        }
        if (screenCnt == null)
        {
            Console.WriteLine("❌ Контур не найден");
            return;
        }
        // Коррекция (упрощённо)
        string outPath = Path.GetFileNameWithoutExtension(input) + "_scanned.jpg";
        Cv2.ImWrite(outPath, src);
        Console.WriteLine($"💾 Сохранён: {outPath}");
        Console.WriteLine($"⏱️ Время: {(DateTime.Now - start).TotalSeconds:F2} сек.");
    }
}
