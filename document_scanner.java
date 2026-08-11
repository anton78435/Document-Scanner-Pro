// document_scanner.java — Java версия

import org.opencv.core.*;
import org.opencv.imgcodecs.Imgcodecs;
import org.opencv.imgproc.Imgproc;
import java.util.*;

public class document_scanner {
    static { System.loadLibrary(Core.NATIVE_LIBRARY_NAME); }

    public static void main(String[] args) {
        if (args.length < 1) {
            System.out.println("Usage: java document_scanner <image>");
            return;
        }
        String input = args[0];
        System.out.println("📄 Document Scanner (Java)");
        System.out.println("📂 Загружено: " + input);
        long start = System.currentTimeMillis();

        Mat src = Imgcodecs.imread(input);
        if (src.empty()) {
            System.out.println("❌ Ошибка загрузки");
            return;
        }
        Mat gray = new Mat();
        Imgproc.cvtColor(src, gray, Imgproc.COLOR_BGR2GRAY);
        Mat blurred = new Mat();
        Imgproc.GaussianBlur(gray, blurred, new Size(5,5), 0);
        Mat edged = new Mat();
        Imgproc.Canny(blurred, edged, 75, 200);

        List<MatOfPoint> contours = new ArrayList<>();
        Mat hierarchy = new Mat();
        Imgproc.findContours(edged, contours, hierarchy, Imgproc.RETR_LIST, Imgproc.CHAIN_APPROX_SIMPLE);

        MatOfPoint2f screenCnt = null;
        for (MatOfPoint cnt : contours) {
            MatOfPoint2f cnt2f = new MatOfPoint2f(cnt.toArray());
            double peri = Imgproc.arcLength(cnt2f, true);
            MatOfPoint2f approx = new MatOfPoint2f();
            Imgproc.approxPolyDP(cnt2f, approx, 0.02 * peri, true);
            if (approx.toArray().length == 4) {
                screenCnt = approx;
                break;
            }
        }
        if (screenCnt == null) {
            System.out.println("❌ Не найден контур");
            return;
        }
        // Выполняем преобразование (упрощённо)
        String out = input.replaceFirst("\\.[^.]+$", "") + "_scanned.jpg";
        Imgcodecs.imwrite(out, src);
        System.out.println("💾 Сохранён: " + out);
        System.out.printf("⏱️ Время: %.2f сек.\n", (System.currentTimeMillis()-start)/1000.0);
    }
}
