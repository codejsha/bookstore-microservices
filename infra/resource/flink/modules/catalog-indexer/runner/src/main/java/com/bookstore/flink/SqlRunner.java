package com.bookstore.flink;

import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.apache.flink.table.api.EnvironmentSettings;
import org.apache.flink.table.api.TableEnvironment;

/**
 * Executes the semicolon-terminated statements of the given SQL script files
 * in order, after substituting ${NAME} placeholders from environment variables.
 */
public final class SqlRunner {

    private static final Pattern ENV_PLACEHOLDER = Pattern.compile("\\$\\{([A-Za-z_][A-Za-z0-9_]*)}");

    public static void main(String[] args) throws Exception {
        if (args.length == 0) {
            throw new IllegalArgumentException("usage: SqlRunner <script.sql> [<script.sql> ...]");
        }
        TableEnvironment tableEnv =
                TableEnvironment.create(EnvironmentSettings.newInstance().inStreamingMode().build());
        for (String file : args) {
            for (String statement : parseStatements(substituteEnv(Files.readString(Path.of(file))))) {
                tableEnv.executeSql(statement);
            }
        }
    }

    static String substituteEnv(String script) {
        Matcher matcher = ENV_PLACEHOLDER.matcher(script);
        StringBuilder result = new StringBuilder();
        while (matcher.find()) {
            String value = System.getenv(matcher.group(1));
            if (value == null) {
                throw new IllegalStateException("environment variable not set: " + matcher.group(1));
            }
            matcher.appendReplacement(result, Matcher.quoteReplacement(value));
        }
        matcher.appendTail(result);
        return result.toString();
    }

    static List<String> parseStatements(String script) {
        List<String> statements = new ArrayList<>();
        StringBuilder current = new StringBuilder();
        for (String line : script.split("\n")) {
            String trimmed = line.trim();
            if (trimmed.isEmpty() || trimmed.startsWith("--")) {
                continue;
            }
            current.append(line).append('\n');
            if (trimmed.endsWith(";")) {
                String statement = current.toString().trim();
                statements.add(statement.substring(0, statement.length() - 1));
                current.setLength(0);
            }
        }
        String rest = current.toString().trim();
        if (!rest.isEmpty()) {
            statements.add(rest);
        }
        return statements;
    }

    private SqlRunner() {}
}
