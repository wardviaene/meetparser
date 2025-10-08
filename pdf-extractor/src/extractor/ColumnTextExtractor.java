package extractor;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;

import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.text.PDFTextStripper;
import org.apache.pdfbox.text.PDFTextStripperByArea;
import org.apache.pdfbox.Loader;

import java.awt.Rectangle;

public class ColumnTextExtractor {
    public static void main(String[] args) throws IOException {
        if (args.length < 2) {
            System.out.println("Usage: java ColumnTextExtractor <input.pdf> <output.txt>");
            return;
        }

        String inputFile = args[0];
        String outputFile = args[1];

            PDDocument document = Loader.loadPDF(new File(inputFile));
            if (document.isEncrypted()) {
                document.setAllSecurityToBeRemoved(true);
            }

            // --- Step 1: Inspect first page text ---
            PDFTextStripper stripper = new PDFTextStripper();
            stripper.setStartPage(1);
            stripper.setEndPage(1);
            String firstPageText = stripper.getText(document);
            boolean isTwoColumn = false;
            boolean isThreeColumn = false;
            String fileType = "";

            if(firstPageText.contains("SwimTopia Meet Maestro")) {
                isTwoColumn = true;
                fileType = "SwimTopia Meet Maestro";
            }

            List<Float> verticalLines = null;
            // --- Step 2: inspect columns
            if(!isTwoColumn && document.getNumberOfPages() > 0) {
                VerticalLineDetector detector = new VerticalLineDetector(document);
                verticalLines = detector.countVerticalLinesOnFirstPage();
                System.out.println("Vertical lines on first page: " + verticalLines.size());
                for (Float x : verticalLines) {
                    System.out.println("Vertical line at x = " + x);
                }
                if(verticalLines != null && verticalLines.size() == 3) {
                    isThreeColumn = true;
                } else if(verticalLines != null && verticalLines.size() == 1) {
                    isTwoColumn = true;
                }
            }

            StringBuilder extracted = new StringBuilder();

            if (isTwoColumn) {
                if(!fileType.equals("")) {
                    extracted.append("FileType: SwimTopia Meet Maestro\n");
                }
                
                
                // --- Step 2a: Two-column extraction ---
                for (PDPage page : document.getPages()) {
                    float width = page.getMediaBox().getWidth();
                    float height = page.getMediaBox().getHeight();

                    float cutoff = width / 2;

                    if(verticalLines.size() == 1) {
                        cutoff = verticalLines.get(0);
                    }

                    PDFTextStripperByArea areaStripper = new PDFTextStripperByArea();
                    areaStripper.setSortByPosition(true);

                    Rectangle left = new Rectangle(0, 0, (int)cutoff, (int) height);
                    Rectangle right = new Rectangle((int)cutoff, 0, (int)cutoff, (int) height);

                    areaStripper.addRegion("left", left);
                    areaStripper.addRegion("right", right);

                    areaStripper.extractRegions(page);

                    extracted.append(areaStripper.getTextForRegion("left")).append("\n");
                    extracted.append(areaStripper.getTextForRegion("right")).append("\n");
                }
            } else if (isThreeColumn && verticalLines != null) {
                extracted.append("FileType: Three column filetype\n");
                // --- Step 2a: Three-column extraction ---
                for (PDPage page : document.getPages()) {
                    int height = (int)page.getMediaBox().getHeight();

                    PDFTextStripperByArea areaStripper = new PDFTextStripperByArea();
                    areaStripper.setSortByPosition(true);


                    float prevX = 0;
                    for (int i = 0; i < verticalLines.size(); i++) {
                        float x = verticalLines.get(i);
                        int rectX = (int) prevX;
                        int rectWidth = (int) (x - prevX);
                        Rectangle rect = new Rectangle(rectX, 0, rectWidth, height);
                        areaStripper.addRegion("column" + i, rect);
                        prevX = x;
                    }

                    areaStripper.extractRegions(page);

                    

                    for (int i = 0; i < verticalLines.size(); i++) {
                        extracted.append(areaStripper.getTextForRegion("column" + i)).append("\n");
                    }
                
                }
            } else {
                // --- Step 2b: Normal extraction ---
                stripper.setStartPage(1);
                stripper.setSortByPosition(true); 
                stripper.setEndPage(document.getNumberOfPages());
                extracted.append(stripper.getText(document));
            }

            // --- Step 3: Save output ---
            Files.write(Paths.get(outputFile), extracted.toString().getBytes());
            System.out.println("Extraction complete. Output written to " + outputFile);
    }
}
