<?php
// document_scanner.php — PHP версия (использует Imagick, но для контуров нужен exec, используем OpenCV через командную строку)

if ($argc < 2) {
    echo "Usage: php document_scanner.php <image>\n";
    exit(1);
}
$input = $argv[1];
echo "📄 Document Scanner (PHP)\n";
echo "📂 Загружено: $input\n";
$start = microtime(true);

// Для простоты вызываем внешний Python скрипт с OpenCV (или используем Imagick для простых операций)
// Но здесь мы реализуем упрощённый вариант с помощью Imagick: обрезка по контуру не поддерживается.
// Используем shell_exec для вызова Python (если он установлен)
$pythonScript = __DIR__ . '/document_scanner.py';
if (!file_exists($pythonScript)) {
    echo "❌ Требуется Python скрипт для реальной обработки.\n";
    exit(1);
}
$out = shell_exec("python $pythonScript $input 2>&1");
echo $out;

$elapsed = microtime(true) - $start;
echo "⏱️ Время: " . number_format($elapsed, 2) . " сек.\n";
?>
