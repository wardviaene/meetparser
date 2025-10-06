package extractor;

import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.pdfparser.PDFStreamParser;
import org.apache.pdfbox.contentstream.operator.Operator;
import org.apache.pdfbox.cos.COSNumber;
import org.apache.pdfbox.cos.COSBase;
import org.apache.pdfbox.contentstream.PDContentStream;
import org.apache.pdfbox.pdmodel.common.PDRectangle;
import org.apache.pdfbox.pdmodel.common.PDStream;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Scanner;
import java.util.Iterator;

public class VerticalLineDetector {
    private final PDDocument document;

    private final float minLineHeight;

    public VerticalLineDetector(PDDocument document) {
        this.document = document;
        this.minLineHeight = 20;
    }

    public VerticalLineDetector(PDDocument document, float minLineHeight) {
        this.document = document;
        this.minLineHeight = minLineHeight;
    }

    /**
     * Counts vertical lines on the first page of the PDF.
     */
    public List<Float> countVerticalLinesOnFirstPage() throws IOException {
        if (document.getNumberOfPages() == 0) {
            return null;
        }

        PDPage page = document.getPage(0);
        return countVerticalLinesOnPage(page);
    }

    public List<Float> countVerticalLinesOnPage(PDPage page) throws IOException {
        List<String> tokens = parseContentStream(page);

        float lastX = 0, lastY = 0;
        boolean hasLastPoint = false;

        Map<Float, Float> verticalSegments = new HashMap<>();

        for (int i = 0; i < tokens.size(); i++) {
            String token = tokens.get(i);

            if (token.equals("m")) { // moveTo
                if (i >= 2) {
                    try {
                        lastX = Float.parseFloat(tokens.get(i - 2));
                        lastY = Float.parseFloat(tokens.get(i - 1));
                        hasLastPoint = true;
                    } catch (NumberFormatException e) {
                        hasLastPoint = false;
                    }
                }
            } else if (token.equals("l") && hasLastPoint) { // lineTo
                if (i >= 2) {
                    try {
                        float x = Float.parseFloat(tokens.get(i - 2));
                        float y = Float.parseFloat(tokens.get(i - 1));

                        float dx = Math.abs(x - lastX);
                        float dy = Math.abs(y - lastY);

                        if (dx < 1.0) {
                            float length = verticalSegments.getOrDefault(x, 0f);
                            length += dy;
                            verticalSegments.put(x, length);
                        }

                        lastX = x;
                        lastY = y;
                    } catch (NumberFormatException e) {
                        // ignore malformed tokens
                    }
                }
            }
        }

        List<Float> verticalLines = new ArrayList<>();
        for (Map.Entry<Float, Float> entry : verticalSegments.entrySet()) {
            if (entry.getValue() >= this.minLineHeight) {
                verticalLines.add(entry.getKey());
            }
        }

        Collections.sort(verticalLines);

        return verticalLines;
    }


    private List<String> parseContentStream(PDPage page) throws IOException {
        List<String> tokens = new ArrayList<>();
        Iterator<PDStream> contentStreams = page.getContentStreams();
        while (contentStreams.hasNext()) {
            PDStream stream = contentStreams.next();
            try (InputStream is = stream.createInputStream()) {
                // Use a scanner to read the content as whitespace-separated tokens
                Scanner scanner = new Scanner(is);
                while (scanner.hasNext()) {
                    tokens.add(scanner.next());
                }
                scanner.close();
            }
        }
        return tokens;
        /*for (PDContentStream contentStream : page.getContentStreams()) {
            PDFStreamParser parser = new PDFStreamParser(contentStream);
            parser.parse();
            // Manually process tokens
            while (parser.hasNext()) {
                Object token = parser.next();
                tokens.add(token);
            }
        }
        return tokens;*/
    }}
