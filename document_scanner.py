

### 1. `document_scanner.py` (Python)

```python
# document_scanner.py — Python версия

import cv2
import numpy as np
import sys
import os
from time import time

def order_points(pts):
    # сортируем точки: верхний-левый, верхний-правый, нижний-правый, нижний-левый
    rect = np.zeros((4, 2), dtype="float32")
    s = pts.sum(axis=1)
    rect[0] = pts[np.argmin(s)]
    rect[2] = pts[np.argmax(s)]
    diff = np.diff(pts, axis=1)
    rect[1] = pts[np.argmin(diff)]
    rect[3] = pts[np.argmax(diff)]
    return rect

def four_point_transform(image, pts):
    rect = order_points(pts)
    (tl, tr, br, bl) = rect
    widthA = np.linalg.norm(br - bl)
    widthB = np.linalg.norm(tr - tl)
    maxWidth = max(int(widthA), int(widthB))
    heightA = np.linalg.norm(tr - br)
    heightB = np.linalg.norm(tl - bl)
    maxHeight = max(int(heightA), int(heightB))
    dst = np.array([
        [0, 0],
        [maxWidth - 1, 0],
        [maxWidth - 1, maxHeight - 1],
        [0, maxHeight - 1]], dtype="float32")
    M = cv2.getPerspectiveTransform(rect, dst)
    warped = cv2.warpPerspective(image, M, (maxWidth, maxHeight))
    return warped

def scan_document(image_path):
    image = cv2.imread(image_path)
    if image is None:
        print("❌ Не удалось загрузить изображение.")
        return None
    orig = image.copy()
    ratio = image.shape[0] / 500.0
    image = cv2.resize(image, (int(image.shape[1] / ratio), 500))
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    gray = cv2.GaussianBlur(gray, (5, 5), 0)
    edged = cv2.Canny(gray, 75, 200)
    contours, _ = cv2.findContours(edged.copy(), cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)
    contours = sorted(contours, key=cv2.contourArea, reverse=True)[:5]
    screenCnt = None
    for c in contours:
        peri = cv2.arcLength(c, True)
        approx = cv2.approxPolyDP(c, 0.02 * peri, True)
        if len(approx) == 4:
            screenCnt = approx
            break
    if screenCnt is None:
        print("❌ Не удалось найти четырёхугольный контур.")
        return None
    warped = four_point_transform(orig, screenCnt.reshape(4, 2) * ratio)
    return warped

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python document_scanner.py <image_path>")
        sys.exit(1)
    input_path = sys.argv[1]
    if not os.path.exists(input_path):
        print("Файл не найден")
        sys.exit(1)
    start = time()
    print(f"📄 Document Scanner (Python)")
    print(f"📂 Загружено: {input_path}")
    result = scan_document(input_path)
    if result is not None:
        out_path = os.path.splitext(input_path)[0] + "_scanned.jpg"
        cv2.imwrite(out_path, result)
        print(f"💾 Результат сохранён: {out_path} (размер {result.shape[1]}x{result.shape[0]})")
        print(f"⏱️ Время: {time()-start:.2f} сек.")
    else:
        sys.exit(1)
